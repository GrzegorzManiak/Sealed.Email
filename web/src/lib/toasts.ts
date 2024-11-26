import { toast } from 'svelte-sonner';



const throw_toast = (
    message: string,
    description: string
) => {
    toast.error(message, {
        dismissable: true,
        description,
        important: true,
        action: { 
            label: 'Dismiss',
            onClick: () => toast.dismiss() 
        }
    });
};



const success_toast = (
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
    throw_toast,
    success_toast
};