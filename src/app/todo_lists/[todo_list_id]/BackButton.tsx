'use client';

import { useRouter } from 'next/navigation';

export default function BackButton() {
    const router = useRouter();

    function goBack() {
        router.push('todo_lists');
    }

    return (
        <div
            className="flex gap-4 justify-start cursor-pointer"
            onClick={goBack}
        >
            <span className="text-xl font-bold">&lt;</span>
            <h1 className="text-xl font-bold"> Back to Task List</h1>
        </div>
    );
}
