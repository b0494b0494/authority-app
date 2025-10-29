'use client';

import { useState, useCallback } from 'react';
import axios from 'axios';
import { useDropzone } from 'react-dropzone';
import { Button, TextField, Typography, Container, Box, Paper, Alert } from '@mui/material';

export default function Home() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [token, setToken] = useState('');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleRegister = async () => {
    try {
      const response = await axios.post('http://localhost:8080/register', {
        username,
        password,
      });
      setMessage(response.data.message);
    } catch (error: any) {
      setMessage(error.response?.data?.error || 'Registration failed');
    }
  };

  const handleLogin = async () => {
    try {
      const response = await axios.post('http://localhost:8080/login', {
        username,
        password,
      });
      setToken(response.data.token);
      setMessage('Login successful!');
    } catch (error: any) {
      setMessage(error.response?.data?.error || 'Login failed');
    }
  };

  const handleProtected = async () => {
    try {
      const response = await axios.get('http://localhost:8080/protected', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setMessage(response.data.message);
    } catch (error: any) {
      setMessage(error.response?.data?.error || 'Access to protected route failed');
    }
  };

  const onDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      setSelectedFile(acceptedFiles[0]);
      setMessage(`Selected file: ${acceptedFiles[0].name}`);
    }
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    maxFiles: 1,
    accept: { 'application/zip': ['.zip'] },
  });

  const handleFileUpload = async () => {
    if (!selectedFile) {
      setMessage('Please select a file to upload.');
      return;
    }
    if (!token) {
      setMessage('Please log in first to upload files.');
      return;
    }

    const formData = new FormData();
    formData.append('file', selectedFile);

    try {
      const response = await axios.post('http://localhost:8080/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          Authorization: `Bearer ${token}`,
        },
      });
      setMessage(response.data.message);
      setSelectedFile(null); // Clear selected file after successful upload
    } catch (error: any) {
      setMessage(error.response?.data?.error || 'File upload failed');
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ my: 4, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Typography component="h1" variant="h3" gutterBottom>
          Authority App
        </Typography>

        <Paper elevation={3} sx={{ p: 4, mt: 3, width: '100%' }}>
          <Typography component="h2" variant="h5" gutterBottom>
            Register / Login
          </Typography>
          <TextField
            margin="normal"
            required
            fullWidth
            id="username"
            label="Username"
            name="username"
            autoComplete="username"
            autoFocus
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            name="password"
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            sx={{ mb: 2 }}
          />
          <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
            <Button
              fullWidth
              variant="contained"
              color="primary"
              onClick={handleRegister}
            >
              Register
            </Button>
            <Button
              fullWidth
              variant="contained"
              color="success"
              onClick={handleLogin}
            >
              Login
            </Button>
          </Box>

          {token && (
            <Box sx={{ mt: 3, textAlign: 'center' }}>
              <Button
                fullWidth
                variant="contained"
                color="secondary"
                onClick={handleProtected}
                sx={{ mb: 2 }}
              >
                Access Protected Route
              </Button>

              <Paper
                {...getRootProps()}
                sx={{
                  border: '2px dashed #cccccc',
                  borderRadius: '4px',
                  p: 4,
                  textAlign: 'center',
                  cursor: 'pointer',
                  bgcolor: isDragActive ? '#fafafa' : 'transparent',
                }}
              >
                <input {...getInputProps()} />
                {
                  selectedFile ?
                    <Typography>{selectedFile.name}</Typography> :
                    <Typography>
                      Drag 'n' drop a .zip file here, or click to select one
                    </Typography>
                }
                <Typography variant="caption" color="textSecondary">
                  (Only .zip files are accepted)
                </Typography>
              </Paper>
              <Button
                fullWidth
                variant="contained"
                color="warning"
                onClick={handleFileUpload}
                sx={{ mt: 2 }}
                disabled={!selectedFile}
              >
                Upload Zip File
              </Button>
            </Box>
          )}

          {message && (
            <Alert severity={message.includes('failed') ? 'error' : 'info'} sx={{ mt: 2 }}>
              {message}
            </Alert>
          )}

          {token && (
            <Typography variant="caption" display="block" sx={{ mt: 2, wordBreak: 'break-all' }}>
              Token: {token}
            </Typography>
          )}
        </Paper>
      </Box>
    </Container>
  );
}
