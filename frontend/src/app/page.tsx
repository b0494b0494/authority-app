'use client';

import { useState } from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import client from '@/lib/grpcClient'; // gRPCクライアントをインポート
import { CreateDirectoryRequest, CreateDirectoryResponse } from '@/grpc/directory_pb'; // 生成されたメッセージをインポート
import { User } from '@/grpc/directory_pb'; // Userメッセージもインポート
import { Status } from '@grpc/grpc-js/build/src/constants'; // gRPCステータスコード

export default function Home() {
  const [directoryName, setDirectoryName] = useState('');
  const [parentId, setParentId] = useState('');
  const [creatorId, setCreatorId] = useState('user-123'); // 仮のユーザーID
  const [responseMessage, setResponseMessage] = useState('');
  const [error, setError] = useState<string | null>(null);

  const handleCreateDirectory = () => {
    setError(null);
    setResponseMessage('');

    const request = new CreateDirectoryRequest();
    request.setName(directoryName);
    if (parentId) {
      request.setParentId(parentId);
    }
    const creator = new User();
    creator.setId(creatorId);
    request.setCreator(creator);

    client.createDirectory(request, (err, res: CreateDirectoryResponse) => {
      if (err) {
        console.error('gRPC Error:', err);
        if (err.code === Status.UNAVAILABLE) {
          setError('gRPCサーバーに接続できません。バックエンドが起動しているか確認してください。');
        } else {
          setError(`ディレクトリ作成に失敗しました: ${err.message}`);
        }
        return;
      }
      setResponseMessage(`ディレクトリ '${res.getDirectory()?.getName()}' が作成されました (ID: ${res.getDirectory()?.getId()})`);
    });
  };

  return (
    <Container component="main" maxWidth="xs">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Typography component="h1" variant="h5">
          ディレクトリ作成
        </Typography>
        <Box component="form" noValidate sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="directoryName"
            label="ディレクトリ名"
            name="directoryName"
            autoFocus
            value={directoryName}
            onChange={(e) => setDirectoryName(e.target.value)}
          />
          <TextField
            margin="normal"
            fullWidth
            id="parentId"
            label="親ディレクトリID (任意)"
            name="parentId"
            value={parentId}
            onChange={(e) => setParentId(e.target.value)}
          />
          <Button
            type="button"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            onClick={handleCreateDirectory}
          >
            ディレクトリ作成
          </Button>
          {responseMessage && (
            <Typography variant="body1" color="primary">
              {responseMessage}
            </Typography>
          )}
          {error && (
            <Typography variant="body1" color="error">
              {error}
            </Typography>
          )}
        </Box>
      </Box>
    </Container>
  );
}
