export const clamp = (number: number, min: number = 0.0, max: number = 1.0) => {
    return Math.max(min, Math.min(number, max))
}
