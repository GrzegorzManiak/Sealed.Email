import { People} from "./generated/people";

const Author = {
    "@type": "Person",
    "name": "Grzegorz Maniak",
    "url": "https://grzegorz.ie/",
    "sameAs": People.GrzegorzManiak.links
};

const Publisher = {
    "@type": "Organization",
    "name": "Grzegorz Maniak",
    "url": "https://grzegorz.ie/",
    "logo": {
        "@type": "ImageObject",
        "url": "https://grzegorz.ie/favicon.png",
        "width": 60,
        "height": 60
    }
};

export {
    People,
    Publisher,
    Author
}