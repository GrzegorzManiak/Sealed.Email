import type { CanvasManager } from "./canvas";



const generate_text_points = (
    text: string, 
    text_height: number,
    canvas_manager: CanvasManager,
) => {
    const balls: Array<{ x: number, y: number }> = [];
    const matrix: Array<Array<number>> = [];

    const width = canvas_manager.draw_ctx.canvas.width,
        height = canvas_manager.draw_ctx.canvas.height;


    canvas_manager.draw_ctx.clearRect(0, 0, width, height);
    canvas_manager.draw_ctx.font = `${text_height}px Arial`;
    canvas_manager.draw_ctx.fillStyle = '#000';
    canvas_manager.draw_ctx.fillText(text, 0, text_height);

        
    // -- get a Uint32 representation of the bitmap:
    const data32 = new Uint32Array(canvas_manager.draw_ctx.getImageData(
        0, 0,
        width, 
        height 
    ).data.buffer);
    


    for(let i = 0; i < data32.length; i++) {
        // -- Check if the pixel is not transparent
        if (data32[i] & 0xff000000) balls.push({
            x: (i % width) | 0,
            y: ((i / width) | 0),
        });
    }


    // -- Convert the balls into a matrix of points wher 1 is a point and 0 is not
    for(let i = 0; i < balls.length; i++) {
        const ball = balls[i];

        if (!matrix[ball.y]) matrix[ball.y] = [];
        matrix[ball.y][ball.x] = 1;
    }


    // -- Fill in the gaps
    let processes: Array<Array<number>> = [];
    let rows = matrix.length;
    let cols = 0;
    for(let i = 0; i < matrix.length; i++) {
        const row = matrix[i];
        if (!row) continue;
        let new_row = [];
        
        for(let j = 0; j < row.length; j++) {
            if (row.length > cols) cols = row.length;
            if (row[j] === 1) new_row.push(1);
            else new_row.push(0);
        }

        processes.push(new_row);
    }

            

    // -- return 2D array of points
    return {
        dots: processes,
        rows, cols
    };
};



export {
    generate_text_points
};