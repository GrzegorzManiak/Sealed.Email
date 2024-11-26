import { throwToast } from '$lib/toasts';



class CanvasManager {
    private readonly _canvas_element: HTMLCanvasElement;
    private readonly _ctx: CanvasRenderingContext2D; 
    private readonly _draw_ctx: OffscreenCanvasRenderingContext2D;

    // -- Framerates
    private _display_frame_rate: number = 1;
    private _last_display_frame_rate: number = 1;
    private _display_frame_rate_counter: number = 1;
    private _display_frame_sum: number = 1;

    private _last_average_frame_measure: number = 1;
    private _averaged_frame_rate_counter: number = 1;
    private _averaged_frame_sum: number = 1;
    private _averaged_frame_rate: number = 1;

    private _current_frame_rate: number = 0;
    private _target_frame_rate: number = 60;
    private _frame_rate_interval: number = 0;
    private _last_frame_time: number = 0;
    private _frame_time: number = 0;

    // -- Size 
    private _width: number = 0;
    private _height: number = 0;

    // -- Debug
    private _started: boolean = false;
    private _debug: boolean = false;
    private _paused: boolean = false;

    // -- Draw callback
    private _callbacks: Array<(
        ctx: CanvasRenderingContext2D,
        draw_ctx: OffscreenCanvasRenderingContext2D,
        width: number,
        height: number,
        canvas_manager: CanvasManager
    ) => void> = [];

    constructor(canvas_element: HTMLCanvasElement) {
        this._canvas_element = canvas_element;

        // -- Onscreen context
        const ctx = canvas_element.getContext('2d');
        if (ctx === null) {
            throwToast('Canvas Manager Error', 'Could not get 2D context from canvas element');
            throw new Error('Could not get 2D context from canvas element');
        }
        this._ctx = ctx;
        
        // -- Offscreen canvas
        const osc = new OffscreenCanvas(1, 1);
        const draw_ctx = osc.getContext('2d');
        if (draw_ctx === null) {
            throwToast('Canvas Manager Error', 'Could not get 2D context from offscreen canvas');
            throw new Error('Could not get 2D context from offscreen canvas');
        }
        this._draw_ctx = draw_ctx;

        // -- Set initial size
        this._on_size_change();

        // -- Set frame rate interval (This is the target time between frames)
        this._frame_rate_interval = 1000 / (this._target_frame_rate * 0.9);
    };



    // -- Methods
    public start(): void {
        this._last_frame_time = performance.now();
        this._frame_time = 0;
        this._current_frame_rate = 0;
        this._paused = false;

        // -- Start only once
        if (!this._started) {
            this._started = true;
            this._loop();

            this._last_average_frame_measure = performance.now();
            this._last_display_frame_rate = performance.now();

            window.addEventListener('resize', this._on_size_change);
            this._canvas_element.addEventListener('resize', this._on_size_change);
            this._on_size_change();
        }
    };  



    public add_callback = (
        callback: (
            ctx: CanvasRenderingContext2D,
            draw_ctx: OffscreenCanvasRenderingContext2D,
            width: number,
            height: number,
            canvas_manager: CanvasManager
        ) => void
    ): void => {
        this._callbacks.push(callback);
    };



    public get_relative_pos(
        e: MouseEvent | TouchEvent
    ): [number, number] {
        try {
            const rect = this._canvas_element.getBoundingClientRect();
            if (e instanceof MouseEvent) return [e.clientX - rect.left, e.clientY - rect.top];
            else return [e.touches[0].clientX - rect.left, e.touches[0].clientY - rect.top];
        }

        catch {
            return [-5000, -5000];
        }
    };



    private _display_frame_time = (): void => {
        this._display_frame_sum += this._averaged_frame_rate;
        this._display_frame_rate_counter++;

        if (performance.now() - this._last_display_frame_rate >= 1000) {
            this._display_frame_rate = this._display_frame_sum / this._display_frame_rate_counter;
            this._display_frame_sum = 1;
            this._display_frame_rate_counter = 1;
            this._last_display_frame_rate = performance.now();
        }
    };



    private _frame_time_keeper = (current_time: number) => {

        // -- Calculate average frame rate every second
        if (current_time - this._last_average_frame_measure >= this._target_frame_rate * 0.6) {
            this._averaged_frame_rate = this._averaged_frame_sum / this._averaged_frame_rate_counter;
            this._display_frame_time();
            this._averaged_frame_sum = 0;
            this._averaged_frame_rate_counter = 0;
            this._last_average_frame_measure = current_time;
            
        }

        // -- Else accumulate
        this._averaged_frame_sum += this._current_frame_rate;
        this._averaged_frame_rate_counter++;
    };



    private async _loop(): Promise<void> {
        while (this._paused) {
            await CanvasManager.sleep(100);
        };

        const current_time = performance.now();
        this._frame_time_keeper(current_time);


        this._frame_time = current_time - this._last_frame_time;
        this._current_frame_rate = 1000 / this._frame_time;

        // -- If enough time has passed, draw
        if (
            this._frame_time >= this._frame_rate_interval && 
            this._averaged_frame_rate <= this._target_frame_rate
        ) {
            this._clear();
            this._draw();
            this._last_frame_time = performance.now();
        }

        else {
            await CanvasManager.sleep(1);
        }

        requestAnimationFrame(() => this._loop());
    };



    private _draw(): void {
        if (this._debug) this._render_framerate();
        this._callbacks.forEach((callback) => {
            callback(this._ctx, this._draw_ctx, this._width, this._height, this);
        });
    };
    


    private _clear(): void {
        this._ctx.save();
        this._ctx.setTransform(1, 0, 0, 1, 0, 0);
        this._ctx.clearRect(0, 0, this._width, this._height);
        this._ctx.restore();

        this._draw_ctx.save();
        this._draw_ctx.setTransform(1, 0, 0, 1, 0, 0);
        this._draw_ctx.clearRect(0, 0, this._width, this._height);
        this._draw_ctx.restore();
    };



    private _render_framerate = (): void => {
        this._ctx.font = '12px Arial';
        this._ctx.fillStyle = 'black';

        // -- Bottom left
        const bottom_left = this._height - 10;

        // -- White text
        this._ctx.fillStyle = 'white';
        this._ctx.fillText(`FPS: ${this._display_frame_rate.toFixed(2)}`, 10, bottom_left);
        this._ctx.fillText(`Width: ${this._width}`, 10, bottom_left - 15);
        this._ctx.fillText(`Height: ${this._height}`, 10, bottom_left - 30);
        this._ctx.fillText(`Frame Time: ${this._frame_time.toFixed(2)}`, 10, bottom_left - 45);
    };



    private _on_size_change = (): void => {
        this._width = this._canvas_element.offsetWidth;
        this._height = this._canvas_element.offsetHeight;
        this._canvas_element.width = this._width;
        this._canvas_element.height = this._height;
        this._draw_ctx.canvas.width = this._width;
        this._draw_ctx.canvas.height = this._height;
    };



    public static sleep = (ms: number): Promise<void> => {
        return new Promise(resolve => setTimeout(resolve, ms));
    };



    get canvas_element(): HTMLCanvasElement { return this._canvas_element; };
    get ctx(): CanvasRenderingContext2D { return this._ctx; };
    get draw_ctx(): OffscreenCanvasRenderingContext2D { return this._draw_ctx; };

    get width(): number { return this._width; };
    get height(): number { return this._height; };

    get averaged_frame_rate(): number { return this._averaged_frame_rate; };
    get current_frame_rate(): number { return this._current_frame_rate; };
    get target_frame_rate(): number { return this._target_frame_rate; };
    get frame_rate_interval(): number { return this._frame_rate_interval; };
    get last_frame_time(): number { return this._last_frame_time; };
    get frame_time(): number { return this._frame_time; };
    get debug(): boolean { return this._debug; };
    get paused(): boolean { return this._paused; };

    set target_frame_rate(value: number) { this._target_frame_rate = value; };
    set debug(value: boolean) { this._debug = value; };
    set paused(value: boolean) { this._paused = value; };
};



export {
    CanvasManager
};