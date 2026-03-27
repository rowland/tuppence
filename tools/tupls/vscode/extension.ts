import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';
import {
    LanguageClient,
    LanguageClientOptions,
    ServerOptions,
    TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
    // Resolve the server binary relative to the extension's install directory.
    // This avoids any dependency on the user's PATH, which is unreliable when
    // VS Code is launched from the Dock or a GUI launcher on macOS.
    const serverCommand = context.asAbsolutePath(path.join('bin', 'tupls'));

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

    const clientOptions: LanguageClientOptions = {
        documentSelector: [{ scheme: 'file', language: 'tuppence' }],
        synchronize: {
            fileEvents: workspace.createFileSystemWatcher('**/*.tup')
        }
    };

    client = new LanguageClient(
        'tuppenceLanguageServer',
        'Tuppence Language Server',
        serverOptions,
        clientOptions
    );

    client.start();
}

export function deactivate(): Thenable<void> | undefined {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
