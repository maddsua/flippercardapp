import type { ResourceVisibility } from "./api_models";

interface ResourceVisibilityOption {
	value: ResourceVisibility;
	label: string;
};

export const resourceVisibilityOptions: ResourceVisibilityOption[] = [
	{
		value: 'PUBLIC',
		label: 'Public; Everyone can see and access it'
	},
	{
		value: 'HIDDEN',
		label: 'Hidden; Users need a link to access'
	},
	{
		value: 'PRIVATE',
		label: 'Private; For team members only'
	},
];
