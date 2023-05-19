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

    try {
        const tagsResponse = await prisma.tags.create({
            data: {
                name: body.name,
                user_id: body.user_id,
            },
        });

        console.log(tagsResponse);

        return NextResponse.json({ ok: true, data: tagsResponse });
    } catch (error) {
        return new NextResponse(JSON.stringify({ error }), { status: 500 });
    }
}
