import { isReactive, isRef, toRaw } from "vue";

//	this basically calls toRaw recuresively
//	note: it doesn't support all the structure-cloneable types like data views and what not
export const toRawDeep = <T>(observed: T): T => {

	const val = isRef<T>(observed) ? observed.value : isReactive(observed) ? toRaw(observed) : observed;

	if (val === null || typeof val !== 'object' || val instanceof Date) {
		return val;
	}

	if (Array.isArray(val)) {
		return val.map(toRawDeep) as T;
	}

	if (val instanceof Set) {
		return new Set(Array.from(val.values()).map(toRawDeep)) as T;
	}

	if (val instanceof Map) {
		return new Map(Array.from(val.entries()).map(([key, val]) => ([key, toRawDeep(val)]))) as T;
	}

	const rawVal: Record<string | symbol, any> = {};
	for (const key of Reflect.ownKeys(val)) {
		rawVal[key] = toRawDeep((val as any)[key]);
	}

	return rawVal as T;
};

export const reactiveSnapshot = <T>(observed: T): T => structuredClone(toRawDeep(observed));
