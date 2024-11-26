<script lang='ts'>
	import { cn } from '$lib/utils.js';
    import { onMount } from 'svelte';
    import type { HTMLAttributes } from 'svelte/elements';
    import { CanvasManager } from './canvas';
    import { throwToast } from '$lib/toasts';
    import { generate_text_points } from './text';
    import {DefaultSettings, type Dots, DotsWritable, Large, Medium, Small} from './index';
    import { render_dots } from './dots';

    // -- Styles
    type $$Props = HTMLAttributes<HTMLCanvasElement>;
	let class_names: $$Props['class'] = undefined;
	export { class_names as class };

    let canvas_element: HTMLCanvasElement;
    let mouse_event: MouseEvent | TouchEvent;

    $: settings = DefaultSettings;

    // -- This is where everything gets launched from
    onMount(() => {
        const screenW = window.innerWidth;

        const sizes = [Small, Medium, Large];
        let size = sizes[0];
        for (let i = 0; i < sizes.length; i++) {
            if (screenW > sizes[i].width) size = sizes[i];
        }

        console.log('Screen size', screenW, size);
        settings.dot_size = size.dot_size;
        settings.dot_spacing = size.dot_spacing;
        settings.text_height = size.text_height;

        // -- Ensure that the canvas element exists
        if (!canvas_element) {
            throwToast('Canvas element not found', 'Canvas element did not mount properly');
            throw new Error('Canvas element not found');
        }

        console.log('Canvas component Mounted');
        const canvas_manager = new CanvasManager(canvas_element);
        canvas_manager.target_frame_rate = 100;
        canvas_manager.debug = false;
        canvas_manager.start();

        let balls = generate_text_points(settings.text, settings.text_height, canvas_manager);
        settings.data = balls;
        console.log('Settings loaded', settings);

        DotsWritable.subscribe((value) => {
            settings = value;
            settings.data = generate_text_points(settings.text, settings.text_height, canvas_manager);
            console.log('Settings updated');
        });


        canvas_manager.add_callback((
            ctx: CanvasRenderingContext2D,
            draw_ctx: OffscreenCanvasRenderingContext2D,
            width: number,
            height: number,
            canvas_manager: CanvasManager
        ) => {

            // -- The size of the text
            const name_width = settings.data.cols * settings.dot_size + (settings.data.cols * settings.dot_spacing);
            const name_height = settings.data.rows * settings.dot_size + (settings.data.rows * settings.dot_spacing);

            // -- Calculate the center of the text
            const x = (ctx.canvas.width / 2) - (name_width / 2),
                y = (ctx.canvas.height / 2) - (name_height / 2);
            
            render_dots(canvas_manager, settings, mouse_event, x, y)
        });
    });

    function OnTouchEnd(e: TouchEvent) {
        e.preventDefault();
        mouse_event = e;
    }

    function OnTouchMove(e: TouchEvent) {
        e.preventDefault();
        mouse_event = e;
    }
</script>



<div class={cn('relative', class_names)}>
    <canvas 
        class='w-full h-full'
        on:touchend={e => OnTouchEnd(e)}
        on:touchmove={e => OnTouchMove(e)}
        on:mousemove={e => mouse_event = e}
        bind:this={canvas_element}>
        
    </canvas>
</div>