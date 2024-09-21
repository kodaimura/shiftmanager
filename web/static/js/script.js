export const getErrorStatus = (error) => {
    const match = error.message.match(/HTTP Status: (\d+)/);
    const status = match? match[1] : "0";
    return parseInt(status);
}

export const handleResponse = (response) => {
    if (!response.ok) {
        throw new Error(`HTTP Status: ${response.status}`);
    } else {
        return response.json();
    }
}

export const handleError = (error) => {
    const status = getErrorStatus(error);
    if (status === 0) {
        console.error(error);
    }
    throw error;
}