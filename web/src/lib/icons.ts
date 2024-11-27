const assets = '/src/assets/attachments/'

const KnownIcons: Map<string, string> = new Map([
    ['pdf', 'pdf.png'],
]);

function GetIcon(name: string, fallback: string = 'file.png'): string {
    name = name.trim().toLowerCase();
    if (name.includes('.')) name = name.split('.').pop() || '';
    return assets + (KnownIcons.get(name) || fallback);
}

export {
    KnownIcons,
    GetIcon
}