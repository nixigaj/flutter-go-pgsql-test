import 'package:adaptive_theme/adaptive_theme.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:test_drive/constants.dart';
import 'dart:io';
import 'package:window_manager/window_manager.dart';
import 'color_schemes.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  if (Platform.isLinux) {
    await windowManager.ensureInitialized();
    // If non-headerbar is used
    //WindowManager.instance.setMinimumSize(const Size(360+70, 640+107));
    WindowManager.instance.setMinimumSize(const Size(360+52, 640+52));
  } else if (Platform.isWindows || Platform.isMacOS) {
    await windowManager.ensureInitialized();
    WindowManager.instance.setMinimumSize(const Size(360, 640));
  }

  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return AdaptiveTheme(
        light: ThemeData(useMaterial3: true, colorScheme: lightColorScheme),
        dark: ThemeData(useMaterial3: true, colorScheme: darkColorScheme),
        initial: AdaptiveThemeMode.system,
        builder: (theme, darkTheme) => MaterialApp(
          title: 'Flutter Demo',
          theme: theme,
          darkTheme: darkTheme,
          home: const MyHomePage(title: 'Flutter Demo Home Page'),
        )

    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  String _status = "Press the wave button to fetch the database resource";
  String _response = "";
  bool _responseVisible = false;

  String truncate(String text, { length = 7, omission = '...' }) {
    if (length >= text.length) {
      return text;
    }
    return text.replaceRange(length, text.length, omission);
  }

  void _fetchDbResource() async { // Use async for HTTP requests
    setState(() {
      _status = "Fetching resource...";
      _responseVisible = false; // Hide response until ready
    });

    try {
      // Get APP_API from environment variables
      var appApi = Platform.environment['APP_API'];

      // Otherwise use default
      appApi ??= defaultAppApiUrl;

      final uri = Uri.parse('$appApi/sql-hello');
      final response = await http.get(uri);

      if (response.statusCode == 200) {
        setState(() {
          _response = truncate(response.body, length: 100);
          _status = "Response:";
          _responseVisible = true; // Show the response now
        });
      } else {
        setState(() {
          _status = "Error: HTTP status code: ${response.statusCode}";
        });
      }
    } catch (error) {
      setState(() {
        _status = "Failed to fetch database resource";
        print(error);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    // This method is rerun every time setState is called, for instance as done
    // by the _incrementCounter method above.
    //
    // The Flutter framework has been optimized to make rerunning build methods
    // fast, so that you can just rebuild anything that needs updating rather
    // than having to individually change instances of widgets.
    return Scaffold(
      appBar: AppBar(
        // TRY THIS: Try changing the color here to a specific color (to
        // Colors.amber, perhaps?) and trigger a hot reload to see the AppBar
        // change color while the other colors stay the same.
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        // Here we take the value from the MyHomePage object that was created by
        // the App.build method, and use it to set our appbar title.
        title: Text(widget.title),
      ),
      body: Center(
        // Center is a layout widget. It takes a single child and positions it
        // in the middle of the parent.
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            // Column is also a layout widget. It takes a list of children and
            // arranges them vertically. By default, it sizes itself to fit its
            // children horizontally, and tries to be as tall as its parent.
            //
            // Column has various properties to control how it sizes itself and
            // how it positions its children. Here we use mainAxisAlignment to
            // center the children vertically; the main axis here is the vertical
            // axis because Columns are vertical (the cross axis would be
            // horizontal).
            //
            // TRY THIS: Invoke "debug painting" (choose the "Toggle Debug Paint"
            // action in the IDE, or press "p" in the console), to see the
            // wireframe for each widget.
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              AnimatedContainer( // Wrap the status in an AnimatedContainer
                duration: const Duration(milliseconds: 200),
                padding: EdgeInsets.only(bottom: _responseVisible ? 20 : 0),
                child: Text(_status),
              ),
              AnimatedSwitcher(
                duration: const Duration(milliseconds: 100),
                transitionBuilder: (Widget child, Animation<double> animation) {
                  return ScaleTransition( // Use ScaleTransition
                    scale: animation,
                    child: child,
                  );
                },
                child: _responseVisible
                    ? Text(
                  _response,
                  key: ValueKey(_response), // Key helps with transition
                  style: Theme.of(context).textTheme.headlineMedium,
                )
                    : const SizedBox.shrink(), // Empty placeholder when not visible
              ),
            ],
          ),
        ),

      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _fetchDbResource,
        tooltip: 'Fetch database resource',
        child: const Icon(Icons.waving_hand_outlined),
      ), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}
