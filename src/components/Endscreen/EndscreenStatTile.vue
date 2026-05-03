<script setup lang="ts">
const props = defineProps<{
	icon?: 'time' | 'target' | 'prize';
}>();
</script>

<template>
	<div class="endscreen-stat-tile">
		<div class="header">
			<div v-if="icon" class="icon" :class="[icon]"></div>
			<span class="title">
				<slot name="title">
					[Title]
				</slot>
			</span>
		</div>
		<div class="content">
			<span class="value">
				<slot>
					[Value]
				</slot>
			</span>
			<span v-if="$slots.after" class="after-value">
				<slot name="after" />
			</span>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.endscreen-stat-tile {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		padding: 0.5rem;
		background-color: var(--app-theme-midnight-glow);
		border-radius: 0.5rem;
		min-width: 0;

		.header {
			display: flex;
			flex-flow: row nowrap;
			gap: 0.5rem;
			align-items: center;
			min-width: 0;

			.icon {
				display: block;
				width: 1rem;
				height: 1rem;
				mask-repeat: no-repeat;
				mask-position: center;
				mask-size: contain;
				flex-shrink: 0;

				&.prize {
					mask-image: url(/src/assets/icons/prize-mask.svg);
					background-color: #FFB900;
				}

				&.time {
					mask-image: url(/src/assets/icons/time-mask.svg);
					background-color: #C27AFF;
				}

				&.target {
					mask-image: url(/src/assets/icons/target-mask.svg);
					background-color: #51A2FF;
				}
			}

			.title {
				font-size: 0.85rem;
				color: var(--app-theme-mysterious-white);
				overflow: hidden;
				white-space: nowrap;
				text-overflow: ellipsis;
			}
		}

		.content {
			display: flex;
			flex-flow: row wrap;
			gap: 0.125rem;
			align-items: center;

			.value {
				font-size: 1.5rem;
				font-weight: 600;
			}

			.after-value {
				color: rgba(255, 255, 255, 0.5);
				font-size: 1.125rem;
				font-weight: 500;
			}
		}
	}
</style>
