<script lang="ts">
    import { goto } from '$app/navigation';
    import { login } from '$lib/api';
    
    let username = '';
    let password = '';
    let error = '';

    async function handleLogin() {
        try {
            const response = await login(username, password);

            if (response.ok) {
                const { token } = await response.json();
                localStorage.setItem('authToken', token);
                await goto('/files');
            } else {
                error = 'Invalid credentials';
            }
        } catch (e) {
            error = 'Login failed';
        }
    }
</script>

<div class="login-container">
    <h1>Login</h1>
    
    {#if error}
        <div class="error">{error}</div>
    {/if}
    
    <form on:submit|preventDefault={handleLogin}>
        <div class="input-group">
            <label for="username">Username</label>
            <input 
                type="text" 
                id="username" 
                bind:value={username} 
                required
            />
        </div>
        
        <div class="input-group">
            <label for="password">Password</label>
            <input 
                type="password" 
                id="password" 
                bind:value={password} 
                required
            />
        </div>
        
        <button type="submit">Login</button>
    </form>

    <div class="signup-link">
        Have an invitation code? <a href="/signup">Sign up here</a>
    </div>
</div>

<style>
    .login-container {
        max-width: 400px;
        margin: 100px auto;
        padding: 2rem;
        border-radius: 8px;
        box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        background: white;
    }

    .input-group {
        margin-bottom: 1rem;
    }

    label {
        display: block;
        margin-bottom: 0.5rem;
    }

    input {
        width: 100%;
        padding: 0.5rem;
        border: 1px solid #ccc;
        border-radius: 4px;
    }

    button {
        width: 100%;
        padding: 0.75rem;
        background: #4a9eff;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }

    button:hover {
        background: #2b7fd9;
    }

    .error {
        color: red;
        margin-bottom: 1rem;
        padding: 0.5rem;
        background: #fee;
        border-radius: 4px;
    }

    .signup-link {
        margin-top: 1rem;
        text-align: center;
        font-size: 0.9rem;
    }

    .signup-link a {
        color: #4a9eff;
        text-decoration: none;
    }

    .signup-link a:hover {
        text-decoration: underline;
    }
</style>