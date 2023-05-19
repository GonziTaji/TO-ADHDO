import { Task } from '@/types';
import { NextResponse } from 'next/server';

export async function POST(request: Request) {
    const body = await request.json();

    console.log(body);

    if (!body.name) {
        return new NextResponse(
            JSON.stringify({ message: 'task must have a name' }),
            { status: 400 }
        );
    }

    const client = await db.connect();

    try {
        const tasksResponse =
            await client.sql`INSERT INTO TASKS (user_id, name) values (${1}, ${
                body.name
            }) returning task_id`;

        console.log(tasksResponse);

        const [{ task_id }] = tasksResponse.rows;

        if (Array.isArray(body.tags) && body.tags.length) {
            const tagsResponse =
                await client.sql`INSERT INTO TAGS (user_id, task_id, name) values
            ${body.tags.map((tag: string) => `(1, ${task_id}, ${tag}),`)}`;
        }

        return JSON.stringify({ ok: true });
    } catch (error) {
        return { error, tasks: [] as Task[] };
    }
}
