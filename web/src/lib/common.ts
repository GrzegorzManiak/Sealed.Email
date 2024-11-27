function RandomHEXColor() {
    const [r, g, b] = [0, 0, 0].map(() => Math.floor(Math.random() * 256).toString(16).padStart(2, '0'));
    return `#${r}${g}${b}`;
}

export { RandomHEXColor };