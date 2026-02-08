import {TelegramAuthProps, useTelegram} from "@/features/auth/model/use-telegram";
import {useRouter} from 'next/navigation'

const TelegramAuthButton = ({
                                buttonSize = 'large',
                                cornerRadius,
                                requestAccess = true
                            }: TelegramAuthProps) => {
    const router = useRouter();

    const {containerRef, isLoading, error} = useTelegram({
        buttonSize,
        cornerRadius,
        requestAccess,
        onSuccess: (data) => {
            console.log('Успешная авторизация:', data);
            router.push('/health/telegram');
        },
        onError: (error) => {
            console.error('Ошибка авторизации:', error);
        }
    });

    return (
        <div className="flex flex-col items-center gap-4">
            <div ref={containerRef} />

            {isLoading && (
                <div className="text-sm text-gray-600">
                    Авторизация...
                </div>
            )}

            {error && (
                <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">
                    Ошибка входа. Попробуйте еще раз.
                </div>
            )}
        </div>
    );
}

export {TelegramAuthButton}