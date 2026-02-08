import {useEffect, useRef} from "react";
import {TelegramAuthError, TelegramAuthResponse, UserTelegram} from "@/entities/user";
import {useMutation} from "@tanstack/react-query";
import {api} from "@/shared";
import {AxiosError} from "axios";

type TelegramAuthProps = {
    buttonSize?: 'large' | 'medium' | 'small';
    cornerRadius?: number;
    requestAccess?: boolean;
    onSuccess?: (data: TelegramAuthResponse) => void;
    onError?: (error: TelegramAuthError) => void;
}

declare global {
    interface Window {
        onTelegramAuth?: (user: UserTelegram) => void;
    }
}

const useTelegram = (
    {
        buttonSize = 'large',
        cornerRadius,
        requestAccess = true,
        onSuccess,
        onError
    }: TelegramAuthProps
) => {
    const containerRef = useRef<HTMLDivElement>(null);

    const {mutate, data, isPending, error} = useMutation<TelegramAuthResponse, AxiosError<TelegramAuthError>, UserTelegram>({
        mutationKey: ['telegram-auth'],
        mutationFn: (user: UserTelegram) => api.post('/api/auth/telegram', user),
        onSuccess: (data) => {
            if (data?.token) {
                localStorage.setItem('authToken', data.token);
            }

            onSuccess?.(data);
        },
        onError: error => {
            console.error('Не получилось авторизоваться через телеграмм: ', error);

            onError?.(error);
        }
    })

    useEffect(() => {
        window.onTelegramAuth = (user: UserTelegram) => {
            mutate(user);
        };

        const script = document.createElement('script');
        script.src = 'https://telegram.org/js/telegram-widget.js?22';
        script.async = true;
        script.setAttribute('data-telegram-login', 'autopost_auth_bot');
        script.setAttribute('data-size', buttonSize);
        script.setAttribute('data-onauth', 'onTelegramAuth(user)');
        
        if (requestAccess) {
            script.setAttribute('data-request-access', 'write');
        }
        
        if (cornerRadius) {
            script.setAttribute('data-radius', cornerRadius.toString());
        } 
        
        containerRef.current?.appendChild(script);
        
        return () => {
            if (containerRef.current) {
                containerRef.current.innerHTML = '';
            }
            delete window.onTelegramAuth;
        };
        
    }, [buttonSize, cornerRadius, mutate, requestAccess]);

    return {
        containerRef,
        isLoading: isPending,
        error,
        data
    }
}

export {useTelegram};
export type {TelegramAuthProps};