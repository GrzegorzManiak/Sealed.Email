import { toast } from 'svelte-sonner';



const throwToast = (
    message: string,
    description: string
) => {
    toast.error(message, {
        dismissable: true,
        description,
        important: true,
    });
};



const successToast = (
    message: string,
    description: string
) => {
    toast.success(message, {
        dismissable: true,
        description,
        important: true,
        action: { 
            label: 'Dismiss',
            onClick: () => toast.dismiss() 
        }
    });
}



export {
    throwToast,
    successToast
};