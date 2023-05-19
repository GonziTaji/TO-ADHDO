import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

async function main() {
    const response = await Promise.all([
        prisma.users.upsert({
            where: { username: 'Yoghurt' },
            update: {},
            create: {
                username: 'Yoghurt',
                Tags: {
                    createMany: {
                        data: [
                            { name: 'gaming' },
                            { name: 'pato' },
                            { name: 'compras' },
                            { name: 'aseo' },
                            { name: 'pieza' },
                            { name: 'trabajo' },
                        ],
                    },
                },
            },
        }),
    ]);

    console.log(response);
}
main()
    .then(async () => {
        await prisma.$disconnect();
    })
    .catch(async (e) => {
        console.error(e);
        await prisma.$disconnect();
        process.exit(1);
    });
