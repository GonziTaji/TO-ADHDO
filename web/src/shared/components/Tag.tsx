import styles from './Tag.module.css'

interface TagProps {
  name: string
  count?: number
  selected?: boolean
  disabled?: boolean
  onClick?: () => void
}

export function Tag({ name, count, selected, disabled, onClick }: TagProps) {
  return (
    <button
      type="button"
      className={`${styles.tag} ${selected ? styles.selected : ''}`}
      disabled={disabled}
      onClick={onClick}
    >
      {name}
      {count !== undefined && <span className={styles.count}>({count})</span>}
    </button>
  )
}
