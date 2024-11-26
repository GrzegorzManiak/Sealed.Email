const RANDOM_TITLES: Array<string> = [
    '- Profesional something',
    '- Will figure it out someday',
    '- I am a developer 🤓',
    '- Um, what?',
    '- I am a human',
    '- I am a robot',
    '- 418 I am a teapot',
];

const TITLE = 'Grzegorz Maniak';

const random_title = () => {
    return TITLE + ' ' + RANDOM_TITLES[Math.floor(Math.random() * RANDOM_TITLES.length)];
};

export {
    random_title,
    TITLE,
    RANDOM_TITLES
};