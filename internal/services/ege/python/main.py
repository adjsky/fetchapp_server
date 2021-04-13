import click
import scripts


@click.group()
def cli():
    pass


@click.command(short_help="solve exam problem")
@click.argument("number", type=int)
@click.option("-f", "--file", "file_", type=click.File(), help="file with data")
@click.option("-t", "--type", "type_", type=int, help="question type")
@click.option("-c", "--char", type=str, help="char to count")
def solve(number, file_, type_, char):
    """
        Solve an exam problem

        NUMBER is a question number
    """
    if scripts.question_implemented(number):
        func = getattr(scripts, "solve_"+str(number))
        func(file_, type_, char)
    else:
        scripts.fatal("Can't solve this question.")


@click.command(short_help="print available questions")
def available():
    """
        Print available questions
    """
    questions_list = scripts.get_available()
    print(", ".join(map(str, questions_list)))


@click.command(short_help="print available question types")
@click.argument("number", type=int)
def types(number):
    """
        Print available question types

        NUMBER is a question number
    """
    scripts.print_available_types(number)


cli.add_command(solve)
cli.add_command(available)
cli.add_command(types)

if __name__ == "__main__":
    cli()
