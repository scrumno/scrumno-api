'use client'

import {api} from "@/shared";
import {useQuery} from "@tanstack/react-query";
import {TelegramAuthButton} from "@/features/auth";

type HealthResponse = {
    message: string,
}

export default function Home() {
    const {data, isLoading} = useQuery({
        queryKey: ['health'],
        queryFn: () => api.get<HealthResponse>('/health/check')
    })

  return (
      <main style={{ padding: '50px', textAlign: 'center', fontSize: '24px' }}>
        <h1>Фронтенд на Next.js</h1>
        <p>Ответ от Go:
            <b> {isLoading ? 'Загрузка...' : data?.data.message}</b>
        </p>
          <TelegramAuthButton buttonSize={'large'} requestAccess={true}/>
      </main>
  )
}