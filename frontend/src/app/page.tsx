'use client'
import { useEffect, useState } from 'react'

export default function Home() {
  const [message, setMessage] = useState('Загрузка...')

  useEffect(() => {
    // Стучимся к нашему Go-серверу
    fetch('/api/health/check')
        .then(res => res.json())
        .then(data => setMessage(data.message))
        .catch(() => setMessage('Ошибка связи с бэкендом'))
  }, [])

  return (
      <main style={{ padding: '50px', textAlign: 'center', fontSize: '24px' }}>
        <h1>Фронтенд на Next.js</h1>
        <p>Ответ от Go: <b>{message}</b></p>
      </main>
  )
}