import Link from 'next/link';
import './globals.css';
import { Inter } from 'next/font/google';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
    title: 'To-adhDo',
    description: '',
};

const links = [
    { route: '/', label: 'Home' },
    { route: '/tasks', label: 'Tasks' },
    { route: '/tags', label: 'Tags' },
];

export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <html lang="en">
            <body className={inter.className + ' h-screen'}>
                <main
                    className="grid overflow-y-hidden h-full"
                    style={{ gridTemplateRows: 'auto 1fr' }}
                >
                    <nav className="flex gap-6 ps-4 py-2 bg-cyan-950 text-white font-bold text-xl shadow-md">
                        {links.map((link, i) => (
                            <Link
                                key={i}
                                className="hover:text-indigo-100 cursor-pointer"
                                href={link.route}
                            >
                                {link.label}
                            </Link>
                        ))}
                    </nav>

                    <section className="overflow-y-auto px-2 py-2">
                        <div className="max-w-xl mx-auto">{children}</div>
                    </section>
                </main>
            </body>
        </html>
    );
}
