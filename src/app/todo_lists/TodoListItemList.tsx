'use client';

import Collapsable from '@/components/Collapsable';
import { Spinner } from '@/components/Spinner';
import { TodoList } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { useState, useTransition } from 'react';

interface TodoListItemListProps {
    todoList: TodoList;
}

export default function TodoListItemList({ todoList }: TodoListItemListProps) {
    const [isPending, startTransition] = useTransition();

    const router = useRouter();

    function navigateToTodoList() {
        startTransition(() => {
            router.push('todo_lists/' + todoList.id);
        });
    }

    return (
        <div
            className={
                'border border-gray-400 p-2 ' +
                (isPending ? 'cursor-wait' : 'cursor-pointer')
            }
            onClick={navigateToTodoList}
        >
            <div className="flex gap-2 items-center">
                <Collapsable vertical={true} collapsed={!isPending}>
                    <div>
                        <Spinner />
                    </div>
                </Collapsable>

                <div>
                    <span className="">{todoList.name}</span>

                    <small className="block">
                        {todoList.task_count}
                        {todoList.task_count === 1 ? 'Task' : 'Tasks'}
                    </small>
                </div>
            </div>
        </div>
    );
}
