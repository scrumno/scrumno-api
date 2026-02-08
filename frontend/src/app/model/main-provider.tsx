'use client'
import {QueryClient} from "@tanstack/query-core";
import {QueryClientProvider} from "@tanstack/react-query";



const MainProvider = ({children}: {children: React.ReactNode}) => {
    const client = new QueryClient()
    return (
        <QueryClientProvider client={client}>
                {children}
        </QueryClientProvider>
    )
}

export {MainProvider}