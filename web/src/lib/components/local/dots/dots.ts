import type { CanvasManager } from './canvas';
import type { Dots } from './index.d';



const render_dots = (
    canvas_manager: CanvasManager, 
    dots: Dots,
    e: MouseEvent | TouchEvent | null = null,
    x_o: number = 0,
    y_o: number = 0
): void => {

    // -- Get the mouse position
    let mouse_pos: [number, number] = [-5000, -5000];
    if (e !== null) mouse_pos = canvas_manager.get_relative_pos(e);
    else mouse_pos = [-5000, -5000];

    canvas_manager.ctx.beginPath();
    canvas_manager.ctx.arc(mouse_pos[0], mouse_pos[1], 10, 0, 2 * Math.PI);
    canvas_manager.ctx.fill();


    // -- Loop through each row
    for (let row = 0; row < dots.data.rows; row++) {

        // -- Loop through each column
        for (let col = 0; col < dots.data.cols; col++) {
            // -- Check if the dot is 0
            if (
                dots.data.dots.length <= row ||
                dots.data.dots[row].length <= col ||
                dots.data.dots[row][col] !== 1
            ) continue;

            // -- Calculate the x and y coordinates
            let x = (col * dots.dot_size) + (col * dots.dot_spacing) + x_o,
                y = (row * dots.dot_size) + (row * dots.dot_spacing) + y_o;

            // -- Calculate the distance
            const dist = get_dist_from_mouse(mouse_pos, x, y);
            let size = dots.dot_size;

            // -- Push the dot away from the mouse
            if (dist < dots.force_size)
                [x, y] = push_dot(mouse_pos, x, y, dist, dots.force_size, dots.force);

            // -- Check if the distance is greater than the max distance
            if (dist < dots.max_dist) {
                // -- Calculate the size of the dot
                const capped_dist = Math.min(dist, dots.max_dist),
                    norm_dist = (capped_dist / dots.max_dist) - 1;

                // -- Calc the size, if 100, then 1, if 0, then 2
                size = dots.dot_size + (norm_dist * dots.dot_size);
            }
            

            // -- Draw the dot
            canvas_manager.ctx.fillStyle = dots.color;
            
            // -- canvas_managerrcle
            canvas_manager.ctx.beginPath();
            canvas_manager.ctx.arc(x, y, size / 2, 0, 2 * Math.PI);
            canvas_manager.ctx.fill();
        }
    }
};



const get_dist_from_mouse = (
    mouse_pos: [number, number],
    x: number,
    y: number,
): number => {

    // -- Calculate the distance
    const dist = Math.sqrt(
        Math.pow(mouse_pos[0] - x, 2) + 
        Math.pow(mouse_pos[1] - y, 2)
    );

    // -- Return the distance
    return dist;
};



const push_dot = (
    mouse_pos: [number, number],
    x: number,
    y: number,
    dist: number,
    max_dist: number = 100,
    force: number = 1
): [number, number] => {

    // -- If the distance is greater than the max distance, then return the original coordinates
    if (dist > max_dist) return [x, y];

    // -- Calculate the angle
    const angle = Math.atan2(mouse_pos[1] - y, mouse_pos[0] - x);

    // -- How much force to apply (the closer the mouse is, the more force to apply)
    force = force * (1 - (dist / max_dist));

    // -- Calculate the new x and y coordinates
    const new_x = x - (Math.cos(angle) * force),
        new_y = y - (Math.sin(angle) * force);

    // -- Return the new coordinates
    return [new_x, new_y];
};



export {
    render_dots,
    get_dist_from_mouse,
    push_dot
};