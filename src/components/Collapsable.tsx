// https://bionicjulia.com/blog/creating-accordion-component-react-typescript-tailwind

import { ReactElement, useRef } from 'react';

interface CollapsableProps {
    collapsed?: boolean;
    children?: ReactElement | ReactElement[] | string;
    className?: string;
    style?: React.CSSProperties;
}

const Collapsable = ({
    children,
    collapsed,
    className = '',
    ...props
}: CollapsableProps) => {
    const childrenParent = useRef<any>(null);

    return (
        <div
            {...props}
            style={{
                height: collapsed
                    ? '0px'
                    : `${childrenParent.current?.scrollHeight}px`,
            }}
            className={
                'overflow-hidden transition-max-height transform duration-300 ease-in-out ' +
                className
            }
        >
            <div ref={childrenParent}>{children}</div>
        </div>
    );
};

export default Collapsable;
