import ArgumentParser
import MyDequeModule

@main
struct PrintStuff: AsyncParsableCommand {
    mutating func run() async throws {
        fputs("My deque colors\n", stdout)
        var colors = MyDeques.colors
        for color in colors {
            fputs("color: \(String(reflecting: color))\n", stdout)
        }
    }
}
