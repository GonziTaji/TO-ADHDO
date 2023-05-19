import { NextResponse } from 'next/server';
import { validateParams } from '@/utils';
import { prisma } from '@/prismaUtils';

export async function POST(request: Request) {
    const body = await request.json();

    const error = validateParams(body, ['name', 'user_id']);

    if (error) {
        return new NextResponse(JSON.stringify({ message: error }), {
            status: 400,
        });
    }

    const tagsIds: number[] = Array.isArray(body.tags) ? body.tags : [];

    try {
        const tasksResponse = await prisma.tasks.create({
            data: {
                name: body.name,
                user_id: body.user_id,
                Tags: {
                    connect: tagsIds.map((id) => ({ id })),
                },
            },
        });

        console.log(tasksResponse);

        return NextResponse.json({ ok: true, data: tasksResponse });
    } catch (error) {
        return new NextResponse(JSON.stringify({ error }), { status: 500 });
    }
}
