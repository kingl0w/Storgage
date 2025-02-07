<script lang="ts">
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';

    onMount(() => {
        const token = localStorage.getItem('authToken');
        if (token) {
            goto('/files');
        } else {
            goto('/login');
        }
    });
</script>

<div>Redirecting...</div>