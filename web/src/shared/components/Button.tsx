import { type ReactNode, type ButtonHTMLAttributes } from 'react'
import styles from './Button.module.css'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger'
  children: ReactNode
}

export function Button({ variant = 'primary', children, className, ...props }: ButtonProps) {
  return (
    <button className={`${styles.button} ${styles[variant]} ${className || ''}`} {...props}>
      {children}
    </button>
  )
}
