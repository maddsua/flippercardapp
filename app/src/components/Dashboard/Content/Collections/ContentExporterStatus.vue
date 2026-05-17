<script setup lang="ts">
const props = defineProps<{
	operation?: string | null;
	progress?: number;
	warn?: string | null;
	error?: string | null;
}>();

const emit = defineEmits<{
	(e: 'cancel'): void;
}>();

</script>

<template>
	<div class="exporter-status">

		<div class="row">

			<div class="operation">
				<template v-if="operation">
					{{ operation }}
				</template>
				<template v-else>
					Operation name
				</template>
			</div>

			<div v-if="typeof progress === 'number'" class="progress">
				<div class="bar" :style="{ width: `${Math.min(progress * 100, 100)}%` }", :class="{ error: !!error }"></div>
			</div>

			<button type="button" class="cancel" @click="emit('cancel')"></button>

		</div>

		<div v-if="warn" class="row">
			<div class="warn-message">
				Warn: {{ warn }}
			</div>
		</div>

		<div v-if="error" class="row">
			<div class="error-message">
				Error: {{ error }}
			</div>
		</div>

	</div>
</template>

<style lang="scss" scoped>
	.exporter-status {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		padding: 0.25rem 0.5rem;
		background-color: var(--app-theme-ghostly-glow);
		border-radius: 0.5rem;
		border: 1px solid var(--app-theme-powder-trail);

		.row {
			display: flex;
			flex-flow: row nowrap;
			gap: 1rem;
			align-items: center;
			justify-content: space-between;
		}

		.operation {
			font-size: 0.75rem;
			flex-shrink: 0;
			font-weight: 600;
		}

		.progress {
			position: relative;
			height: 3px;
			flex-grow: 1;
			background-color: var(--app-theme-powder-trail);
			border-radius: 0.25rem;
			overflow: hidden;

			.bar {
				position: absolute;
				left: 0;
				top: 0;
				height: 100%;
				background-color: var(--app-theme-deep-lavender);
				transition: all 150ms ease;

				&.error {
					background-color: var(--app-theme-blood-red);
				}
			}
		}

		button.cancel {
			position: relative;
			display: block;
			border: none;
			outline: none;
			background: unset;
			padding: 0;
			border-radius: 0.25rem;
			
			&::after {
				content: "";
				display: block;
				width: 1rem;
				height: 1rem;	
				background-color: var(--app-theme-blood-red);
				mask-type: alpha;
				mask-position: center;
				mask-repeat: no-repeat;
				mask-size: contain;
				mask-image: url(/src/assets/icons/cross-cut-mask.svg);
			}

			&:hover {
				cursor: pointer;
				background-color: var(--app-theme-blood-red);

				&::after {
					background-color: var(--app-theme-snow-white);
				}
			}
		}

		.warn-message, .error-message {
			display: block;
			padding: 0.25rem 1rem;
			font-size: 0.75rem;
			font-weight: 600;
			flex-grow: 1;
			min-width: 0;
			overflow: hidden;
			text-overflow: ellipsis;
			white-space: nowrap;
			border-radius: 0.25rem;
		}

		.warn-message {
			color: var(--app-theme-snow-white);
			background-color: var(--app-theme-spooky-orange);
		}

		.error-message {
			color: var(--app-theme-snow-white);
			background-color: var(--app-theme-blood-red);
		}
	}
</style>
