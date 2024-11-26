// -- X,Y Coordinates
export type Dot = [number, number];

export interface Dots {
    data: {
        dots: Array<Array<number>>,
        rows: number,
        cols: number
    };

    dot_size: number;
    dot_spacing: number;

    max_dist: number;
    force: number;
    force_size: number;

    color: string;
}

export type DotsSettings = {
    text: string;
    text_height: number;
} & Dots;

export type FrameCallback = {
    func: () => void;
    remove: () => void;
};

export type LogType = 'INFO' | 'WARN' | 'ERROR' | 'DEBUG';

// -- The medication units will store data about how much medication is in the
// is needed per part of body, FTU (fingertip unit) differs from part of body
// to part of body, so we need to store that data somewhere.

export type Medication = {
    id: string;
    name: string;
    none?: boolean; // -- Default when no medication is selected
};

export interface DrawInstructions {
    fill?: boolean;
    fill_color?: string;

    stroke?: boolean;
    stroke_color?: string;

    line_width?: number;
    line_dash?: Array<number>;
    
    line_cap?: CanvasLineCap;
    line_join?: CanvasLineJoin;

    scale?: boolean;
    draw_outline?: boolean;
}