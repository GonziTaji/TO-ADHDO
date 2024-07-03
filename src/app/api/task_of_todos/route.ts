import { NextResponse } from 'next/server';
import { prisma } from '@/prismaUtils';
import { validateParams } from '@/utils';

export async function POST(request: Request) {
    const body = await request.json();

    const error = validateParams(body, ['task_id', 'tag_id']);

    if (error) {
        return new NextResponse(JSON.stringify({ message: error }), {
            status: 400,
        });
    }

    try {
        const tasksResponse = await prisma.todoLists.update({
            where: { id: body.task_id },
            data: {
                Tasks: {
                    connect: {
                        task_id_todolist_id: body.tag_id,
                    },
                },
            },
        });

        console.log(tasksResponse);

        return NextResponse.json({ ok: true, data: tasksResponse });
    } catch (error: any) {
        return new NextResponse(JSON.stringify({ error: error.toString() }), {
            status: 500,
        });
    }
}

export async function DELETE(request: Request) {
    const url = new URL(request.url);

    const params = {
        task_id: parseInt(url.searchParams.get('task_id') || ''),
        tag_id: parseInt(url.searchParams.get('tag_id') || ''),
    };

    console.log(params);

    const error = validateParams(params, ['task_id', 'tag_id']);

    if (error) {
        return new NextResponse(JSON.stringify({ message: error }), {
            status: 400,
        });
    }

    try {
        const tasksResponse = await prisma.tasks.update({
            where: { id: params.task_id },
            data: {
                Tags: {
                    disconnect: {
                        id: params.tag_id,
                    },
                },
            },
        });

        console.log(tasksResponse);

        return NextResponse.json({ ok: true, data: tasksResponse });
    } catch (error: any) {
        return new NextResponse(JSON.stringify({ error: error.toString() }), {
            status: 500,
        });
    }
}
