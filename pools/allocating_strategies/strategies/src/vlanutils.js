
export function rangeCapacity(vlanRange) {
    return vlanRange.to - vlanRange.from + 1
}

export function rangeToStr(range) {
    return `[${range.from}-${range.to}]`
}

export function freeCapacity(parentRange, utilisedCapacity) {
    return rangeCapacity(parentRange) - utilisedCapacity
}
