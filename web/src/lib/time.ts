function PrettyPrintTime(time: Date): string {
    time = new Date(time);
    // -- If the time is today, return the time
    if (time.toDateString() === new Date().toDateString()) {
        // HH:MM AM/PM
        return time.toLocaleTimeString('en-US', {hour: '2-digit', minute: '2-digit'});
    }

    // -- Else, DAY MMM/DD
    return time.toLocaleDateString('en-US', {weekday: 'short', month: 'short', day: '2-digit'});
}

export {
    PrettyPrintTime
};