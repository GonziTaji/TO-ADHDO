'use client';

import { Spinner } from '@/components/Spinner';
import { useRouter } from 'next/navigation';
import { ChangeEvent, useState, useTransition } from 'react';

export default function NewTodoListForm() {
    const [listName, setListName] = useState('');

    const [isCreating, setIsCreating] = useState(false);
    const [isNavigating, startTransition] = useTransition();

    const isMutating = isCreating || isNavigating;

    const router = useRouter();

    async function createTodoList() {
        setIsCreating(true);

        const response = await fetch('/api/todo_lists', {
            method: 'POST',
            body: JSON.stringify({ user_id: 1, name: listName }),
        });

        const body = await response.json();

        if (body.data && body.data.id) {
            startTransition(() => {
                setIsCreating(false);
                setListName('');

                router.push('/todo_lists/' + body.data.id);
            });
        } else {
            console.error(
                'something happened. No id in response body',
                response
            );
        }
    }

    function inputOnChange(ev: ChangeEvent<HTMLInputElement>) {
        setListName(ev.currentTarget.value);
    }

    return (
        <div className="flex gap-2">
            <input
                type="text"
                value={listName}
                onInput={inputOnChange}
                disabled={isMutating}
                className="ps-2"
                placeholder="Weekend Chores"
            />
            <button
                type="button"
                className="inline-flex whitespace-nowrap items-center gap-3 border rounded cursor-pointer disabled:cursor-not-allowed disabled:bg-green-100 disabled:text-gray-400 border-gray-400 bg-green-400 px-2"
                onClick={createTodoList}
                disabled={isMutating}
            >
                Create List
                {isMutating && <Spinner />}
            </button>
        </div>
    );
}
