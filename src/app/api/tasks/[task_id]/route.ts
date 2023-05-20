import { NextResponse } from 'next/server';
import { prisma } from '@/prismaUtils';
import { validateParams } from '@/utils';

interface RouteParams {
    params: {
        task_id: string;
    };
}

export async function DELETE(request: Request, { params }: RouteParams) {
    const task_id = parseInt(params.task_id);

    try {
        const tasksResponse = await prisma.tasks.delete({
            where: { id: task_id },
        });

        console.log(tasksResponse);

        return NextResponse.json({ ok: true, data: tasksResponse });
    } catch (error: any) {
        return new NextResponse(JSON.stringify({ error: error.toString() }), {
            status: 500,
        });
    }
}

export async function PUT(request: Request, { params }: RouteParams) {
    const task_id = parseInt(params.task_id);

    const body = await request.json();

    if (!Object.keys(body).length) {
        return new NextResponse(JSON.stringify({ error: 'No changes' }), {
            status: 304,
        });
    }

    try {
        const tasksResponse = await prisma.tasks.update({
            where: { id: task_id },
            data: body,
        });

        console.log(tasksResponse);

        return NextResponse.json({ ok: true, data: tasksResponse });
    } catch (error: any) {
        return new NextResponse(JSON.stringify({ error: error.toString() }), {
            status: 500,
        });
    }
}
