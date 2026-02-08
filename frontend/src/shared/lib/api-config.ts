import axios from "axios";

const BASE_URL: string = process.env.NODE_ENV === 'development'
    ? 'http://localhost:8080/api/v1'
    : '/api/v1';

const api = axios.create({
    baseURL: BASE_URL
});

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');

    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }

    config.headers['Content-Type'] = 'application/json';
    config.headers['X-Requested-With'] = 'XMLHttpRequest';

    return config;
});

api.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        if (error.response?.status === 401) {
            console.error("Пользователь не авторизован");
        }
        return Promise.reject(error);
    }
);

export {BASE_URL, api};