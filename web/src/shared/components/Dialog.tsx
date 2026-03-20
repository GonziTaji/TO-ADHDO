import { type ReactNode, useEffect, useRef, type DialogHTMLAttributes } from 'react'
import styles from './Dialog.module.css'

interface DialogProps extends DialogHTMLAttributes<HTMLDialogElement> {
  isOpen: boolean
  onClose: () => void
  title?: string
  children: ReactNode
}

export function Dialog({ isOpen, onClose, title, children, className, ...props }: DialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null)

  useEffect(() => {
    const dialog = dialogRef.current
    if (!dialog) return

    if (isOpen) {
      dialog.showModal()
    } else {
      dialog.close()
    }
  }, [isOpen])

  useEffect(() => {
    const dialog = dialogRef.current
    if (!dialog) return

    const handleClose = () => onClose()
    dialog.addEventListener('close', handleClose)
    return () => dialog.removeEventListener('close', handleClose)
  }, [onClose])

  return (
    <dialog ref={dialogRef} className={`${styles.dialog} ${className || ''}`} {...props}>
      {title && <header><h3>{title}</h3></header>}
      <div className={styles.content}>{children}</div>
      <form method="dialog">
        <button type="submit" onClick={onClose}>Close</button>
      </form>
    </dialog>
  )
}
