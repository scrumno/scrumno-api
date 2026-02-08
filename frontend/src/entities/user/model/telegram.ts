type UserTelegram = {
    id: number;
    first_name: string;
    last_name?: string;
    username: string;
    photo_url?: string;
    auth_date: number;
    hash: string;
}

type TelegramAuthResponse = {
    token: string;
    user: {
        id: number;
        first_name: string;
        last_name?: string;
        username?: string;
        photo_url?: string;
    };
}

type TelegramAuthError = {
    message: string;
    code?: string;
}

export type {UserTelegram, TelegramAuthResponse, TelegramAuthError}