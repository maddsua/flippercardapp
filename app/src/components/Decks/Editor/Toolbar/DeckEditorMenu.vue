<script setup lang="ts">
import { isContainerOutsideClick } from '@/dom';
import { onMounted, onUnmounted, ref } from 'vue';

const props = defineProps<{
	label: string;
}>();

const listShown = ref(false);
const containerRef = ref<HTMLElement | null>(null);

const handleOutsideClick = (event: MouseEvent) => {
	if (listShown.value && isContainerOutsideClick(containerRef.value, event.target)) {
		listShown.value = false;
	}
};

onMounted(() => {
	window.addEventListener('click', handleOutsideClick);
});

onUnmounted(() => {
	window.removeEventListener('click', handleOutsideClick);
});

</script>

<template>
	<div class="editor-menu" ref="containerRef">
		<button type="button" class="label" @click="listShown = !listShown">
			{{ label }}
		</button>
		<div v-if="listShown" class="list-anchor" @click="listShown = false">
			<div class="menu-list">
				<slot>
					<span class="placeholder">
						[Entries]
					</span>
				</slot>
			</div>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.editor-menu {
		position: relative;

		button.label {
			display: block;
			background: unset;
			border: none;
			outline: none;
			border-radius: 0.125rem;
			color: var(--app-theme-kinda-white);
			padding: 0.125rem 0.25rem;

			&:hover {
				cursor: pointer;
				background-color: var(--app-theme-powder-trail);
			}
		}

		.list-anchor {
			position: relative;
			bottom: 0;
			left: 0;

			.menu-list {
				position: absolute;
				top: 0.5rem;
				left: 0;
				width: 100%;
				z-index: 50;
				display: flex;
				flex-direction: column;
				width: 10rem;
				gap: 0.125rem;
				background-color: var(--app-theme-midnight);
				padding: 0.25rem 0;
				border-radius: 0.25rem;
				overflow: hidden;

				span.placeholder {
					color: var(--app-theme-mysterious-white);
					font-size: 0.75rem;
				}
			}
		}
	}
</style>
