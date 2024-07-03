// https://bionicjulia.com/blog/creating-accordion-component-react-typescript-tailwind

import { ReactElement, useRef } from 'react';

interface CollapsableProps {
    collapsed?: boolean;
    children?: ReactElement | ReactElement[] | string;
    className?: string;
    vertical?: boolean;
    style?: React.CSSProperties;
}

const Collapsable = ({
    children,
    collapsed,
    className = '',
    vertical = false,
    ...props
}: CollapsableProps) => {
    const childrenParent = useRef<any>(null);

    let htmlProp = 'height';
    let parentHtmlProp = 'scrollHeight';

    if (vertical) {
        htmlProp = 'width';
        parentHtmlProp = 'scrollWidth';
    }

    return (
        <div
            {...props}
            style={{
                [htmlProp]: collapsed
                    ? '0px'
                    : `${
                          childrenParent.current
                              ? childrenParent.current[parentHtmlProp]
                              : '0'
                      }px`,
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
