function TrimText(text: string, maxLength: number, suffix = "..."): string {
    text = text.trim();
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + suffix;
}

export {
    TrimText
}