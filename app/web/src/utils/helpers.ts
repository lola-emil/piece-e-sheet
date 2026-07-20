

export const formatCurrency = (val: number) => {
    return new Intl.NumberFormat('en-PH', { style: 'currency', currency: 'PHP' }).format(val);
};

export const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString();
};