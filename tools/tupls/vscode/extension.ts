import { workspace, ExtensionContext } from 'vscode';
import {
    LanguageClient,
    LanguageClientOptions,
    ServerOptions,
    TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(_context: ExtensionContext) {
    // The server is implemented in Go
    const serverCommand = 'tupls';

    // If the extension is launched in debug mode then the debug server options are used
    // Otherwise the run options are used
    const serverOptions: ServerOptions = {
        run: {
            command: serverCommand,
            transport: TransportKind.stdio
        },
        debug: {
            command: serverCommand,
            transport: TransportKind.stdio
        }
    };

    // Options to control the language client
    const clientOptions: LanguageClientOptions = {
        // Register the server for Tuppence documents
        documentSelector: [{ scheme: 'file', language: 'tuppence' }],
        synchronize: {
            // Notify the server about file changes to .tup files contained in the workspace
            fileEvents: workspace.createFileSystemWatcher('**/*.tup')
        }
    };

    // Create the language client and start the client.
    client = new LanguageClient(
        'tuppenceLanguageServer',
        'Tuppence Language Server',
        serverOptions,
        clientOptions
    );

    // Start the client. This will also launch the server
    client.start();
}

export function deactivate(): Thenable<void> | undefined {
    if (!client) {
        return undefined;
    }
    return client.stop();
} 