interface PriceDisplayProps {
  price: number
  referencePrice?: number
}

function formatCLP(value: number): string {
  return new Intl.NumberFormat('es-CL', {
    style: 'currency',
    currency: 'CLP',
    minimumFractionDigits: 0
  }).format(value)
}

export function PriceDisplay({ price, referencePrice }: PriceDisplayProps) {
  return (
    <div>
      <span>{formatCLP(price)}</span>
      {referencePrice !== undefined && referencePrice > 0 && (
        <small>
          Ref: {formatCLP(referencePrice)}
        </small>
      )}
    </div>
  )
}
