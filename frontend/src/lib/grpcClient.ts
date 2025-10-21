import { DirectoryServiceClient } from '../grpc/directory_grpc_pb';
import * as grpc from '@grpc/grpc-js';

// バックエンドのgRPCサーバーのアドレス
// Docker Composeでバックエンドが起動している場合、フロントエンドから見たホストはlocalhost
// ポートはバックエンドのgRPCサーバーがリッスンしているポート
const GRPC_SERVER_ADDRESS = process.env.NEXT_PUBLIC_GRPC_SERVER_ADDRESS || 'localhost:50051'; // 例: 50051はGoバックエンドのデフォルトポート

const client = new DirectoryServiceClient(
  GRPC_SERVER_ADDRESS,
  grpc.credentials.createInsecure() // 開発環境ではInsecureでOK
);

export default client;
