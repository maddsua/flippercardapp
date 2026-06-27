<script setup lang="ts">
import EditorModal from '../EditorModal.vue';

interface Shortcut {
	title: string;
	keys: string[];
};

const props = defineProps<{
	shortcuts?: Shortcut[];
}>();

const emit = defineEmits<{
	(e: 'done'): void;
}>();

</script>

<template>
	<EditorModal title="Shortcuts" variant="compact" @close="emit('done')">

		<div class="shortcut-list">

			<div v-for="item of shortcuts" class="keyboard-shortcut">

				<div class="title">
					{{ item.title }}
				</div>

				<div class="key-list">

					<template v-for="(key, idx) of item.keys">

						<div class="key">
							{{ key }}
						</div>

						<div v-if="idx < item.keys.length - 1" class="plus">+</div>

					</template>
				</div>
			</div>

		</div>

	</EditorModal>
</template>

<style lang="scss" scoped>

	.shortcut-list {
		display: flex;
		flex-direction: column;

		.keyboard-shortcut {
			display: flex;
			flex-flow: row nowrap;
			gap: 1rem;
			padding: 0.5rem 1rem;
			background-color: rgba(255, 255, 255, 0.05);
			transition: all 75ms linear;

			&:nth-child(2n) {
				background-color: rgba(255, 255, 255, 0.1);
			}

			&:hover {
				background-color: rgba(255, 255, 255, 0.15);
			}

			.title {
				font-size: 0.85rem;
				font-weight: 300;
				color: var(--app-theme-kinda-white);
			}

			.key-list {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				justify-content: end;
				flex-grow: 1;
				gap: 0.25rem;

				.key {
					font-size: 0.7rem;
					font-weight: 500;
					border-radius: 0.25rem;
					padding: 0.125rem 0.5rem;
					background-color: rgba(255, 255, 255, 0.1);
				}

				.plus {
					font-size: 0.65rem;
					font-weight: 300;
					color: var(--app-theme-kinda-white);
				}
			}
		}
	}

</style>
